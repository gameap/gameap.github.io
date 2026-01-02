/**
 * Custom Table of Contents Generator
 * No external dependencies
 */
(function() {
  'use strict';

  var TOC = {
    config: {
      contentSelector: '#content',
      tocSelector: '#toc',
      headingSelector: 'h2, h3, h4, h5, h6',
      scrollOffset: 100
    },

    headings: [],
    currentActiveId: null,

    init: function() {
      var tocContainer = document.querySelector(this.config.tocSelector);
      var contentContainer = document.querySelector(this.config.contentSelector);

      if (!tocContainer || !contentContainer) {
        return;
      }

      this.headings = contentContainer.querySelectorAll(this.config.headingSelector);

      if (this.headings.length === 0) {
        tocContainer.style.display = 'none';
        return;
      }

      this.ensureHeadingIds(this.headings);

      var tocHTML = this.buildTocHTML(this.headings);
      tocContainer.innerHTML = tocHTML;

      this.initScrollSpy();
      this.setupSmoothScroll(tocContainer);
    },

    ensureHeadingIds: function(headings) {
      var usedIds = {};

      for (var i = 0; i < headings.length; i++) {
        var heading = headings[i];

        if (!heading.id) {
          var baseId = heading.textContent
            .trim()
            .toLowerCase()
            .replace(/[^\w\s-]/g, '')
            .replace(/\s+/g, '-')
            .replace(/-+/g, '-')
            .substring(0, 50);

          var id = baseId;
          var counter = 1;
          while (usedIds[id]) {
            id = baseId + '-' + counter;
            counter++;
          }

          heading.id = id;
        }

        usedIds[heading.id] = true;
      }
    },

    buildTocHTML: function(headings) {
      var items = [];

      for (var i = 0; i < headings.length; i++) {
        var heading = headings[i];
        var level = parseInt(heading.tagName.charAt(1), 10);
        var text = heading.textContent.trim();
        var id = heading.id;

        items.push({
          level: level,
          text: text,
          id: id
        });
      }

      return this.generateNestedList(items);
    },

    generateNestedList: function(items) {
      if (items.length === 0) {
        return '';
      }

      var baseLevel = items[0].level;
      var html = '<nav class="toc-nav"><ul class="toc-list">';
      var currentLevel = baseLevel;

      for (var i = 0; i < items.length; i++) {
        var item = items[i];
        var level = item.level;

        if (level > currentLevel) {
          for (var j = currentLevel; j < level; j++) {
            html += '<ul class="toc-list-nested">';
          }
        } else if (level < currentLevel) {
          for (var k = currentLevel; k > level; k--) {
            html += '</li></ul>';
          }
          html += '</li>';
        } else if (i > 0) {
          html += '</li>';
        }

        currentLevel = level;

        html += '<li class="toc-item toc-h' + level + '">';
        html += '<a class="toc-link" href="#' + item.id + '">' + this.escapeHTML(item.text) + '</a>';
      }

      for (var m = currentLevel; m >= baseLevel; m--) {
        html += '</li>';
        if (m > baseLevel) {
          html += '</ul>';
        }
      }

      html += '</ul></nav>';

      return html;
    },

    escapeHTML: function(text) {
      var div = document.createElement('div');
      div.textContent = text;
      return div.innerHTML;
    },

    initScrollSpy: function() {
      var self = this;
      var contentContainer = document.querySelector(this.config.contentSelector);
      var headings = contentContainer.querySelectorAll(this.config.headingSelector);

      if (!('IntersectionObserver' in window)) {
        // Fallback for older browsers
        this.initScrollSpyFallback();
        return;
      }

      var observer = new IntersectionObserver(function(entries) {
        entries.forEach(function(entry) {
          if (entry.isIntersecting) {
            self.setActiveLink(entry.target.id);
          }
        });
      }, {
        rootMargin: '-80px 0px -70% 0px',
        threshold: 0
      });

      headings.forEach(function(heading) {
        observer.observe(heading);
      });
    },

    initScrollSpyFallback: function() {
      var self = this;
      var ticking = false;

      window.addEventListener('scroll', function() {
        if (!ticking) {
          window.requestAnimationFrame(function() {
            self.updateActiveOnScroll();
            ticking = false;
          });
          ticking = true;
        }
      });

      this.updateActiveOnScroll();
    },

    updateActiveOnScroll: function() {
      var scrollPos = window.pageYOffset + this.config.scrollOffset;
      var activeId = null;

      for (var i = 0; i < this.headings.length; i++) {
        var heading = this.headings[i];
        if (heading.offsetTop <= scrollPos) {
          activeId = heading.id;
        }
      }

      if (activeId) {
        this.setActiveLink(activeId);
      }
    },

    setActiveLink: function(id) {
      if (this.currentActiveId === id) {
        return;
      }

      var tocContainer = document.querySelector(this.config.tocSelector);
      var links = tocContainer.querySelectorAll('.toc-link');

      for (var i = 0; i < links.length; i++) {
        links[i].classList.remove('active');
      }

      var activeLink = tocContainer.querySelector('.toc-link[href="#' + id + '"]');
      if (activeLink) {
        activeLink.classList.add('active');
        this.currentActiveId = id;
      }
    },

    setupSmoothScroll: function(tocContainer) {
      var offset = this.config.scrollOffset;

      tocContainer.addEventListener('click', function(e) {
        if (e.target.tagName === 'A' && e.target.hash) {
          e.preventDefault();

          var targetId = e.target.hash.substring(1);
          var targetElement = document.getElementById(targetId);

          if (targetElement) {
            var targetPosition = targetElement.getBoundingClientRect().top + window.pageYOffset - offset;

            window.scrollTo({
              top: targetPosition,
              behavior: 'smooth'
            });

            history.pushState(null, null, e.target.hash);
          }
        }
      });
    }
  };

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', function() {
      TOC.init();
    });
  } else {
    TOC.init();
  }

  window.TOC = TOC;
})();
