// Wrap entire application in IIFE to prevent duplicate execution and global scope pollution
(function() {
    'use strict';
    
    // Prevent multiple initializations - single source of truth
    if (window.lacpaAppInitialized) {
        console.log('[LACPA] App already initialized, preventing duplicate execution');
        return;
    }
    window.lacpaAppInitialized = true;
    console.log('[LACPA] Initializing app...');

// Custom HTMX event handlers and utilities

// Handle HTMX responses for items list
document.body.addEventListener('htmx:afterSwap', function(event) {
    if (event.detail.target.id === 'items-list') {
        const response = event.detail.xhr.responseText;
        
        try {
            const items = JSON.parse(response);
            
            if (Array.isArray(items) && items.length === 0) {
                event.detail.target.innerHTML = `
                    <div class="text-center py-8 text-gray-500">
                        No items found. Create your first item above!
                    </div>
                `;
            } else if (Array.isArray(items)) {
                event.detail.target.innerHTML = items.map(item => `
                    <div class="border border-gray-200 rounded-lg p-4 hover:shadow-md transition duration-200" id="item-${item.id}">
                        <div class="flex justify-between items-start">
                            <div class="flex-1">
                                <h3 class="text-lg font-semibold text-gray-800 mb-2">${escapeHtml(item.name)}</h3>
                                <p class="text-gray-600 mb-2">${escapeHtml(item.description)}</p>
                                <p class="text-sm text-gray-400">Created: ${new Date(item.created_at).toLocaleString()}</p>
                            </div>
                            <button hx-delete="/api/items/${item.id}"
                                    hx-confirm="Are you sure you want to delete this item?"
                                    hx-target="#item-${item.id}"
                                    hx-swap="outerHTML swap:1s"
                                    class="ml-4 bg-red-500 text-white py-1 px-3 rounded hover:bg-red-600 transition duration-200">
                                Delete
                            </button>
                        </div>
                    </div>
                `).join('');
            }
        } catch (e) {
            console.error('Error parsing response:', e);
        }
    }
});

// Handle successful item creation
document.body.addEventListener('htmx:afterRequest', function(event) {
    if (event.detail.successful && event.detail.xhr.status === 201) {
        // Reload the items list after creating a new item
        const itemsList = document.getElementById('items-list');
        if (itemsList) {
            htmx.trigger(itemsList, 'htmx:trigger');
        }
    }
});

// Utility function to escape HTML
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// ========================================
// NAVIGATION JAVASCRIPT & ANIMATIONS
// ========================================

// Load correct content based on current URL path
function loadPageFromURL() {
    const path = window.location.pathname;
    const search = window.location.search; // Get query parameters like ?page=2
    const mainDiv = document.getElementById('main-div');
    
    if (!mainDiv) {
        console.error('main-div not found');
        return;
    }
    
    // Ensure HTMX is loaded before using it
    if (typeof htmx === 'undefined') {
        console.error('HTMX not loaded, retrying...');
        setTimeout(loadPageFromURL, 100);
        return;
    }
    
    // Map URL paths to API endpoints
    const routeMap = {
        '/': 'http://localhost:3000/api/main/landing/',
        '/discover': 'http://localhost:3000/main/discover',
        '/membership': 'http://localhost:3000/membership',
        '/members/individuals': 'http://localhost:3000/membership', // Legacy support
        '/membership/apply-now': 'http://localhost:3000/membership/apply-now',
        '/membership/firms': 'http://localhost:3000/membership/firms',
        '/events': 'http://localhost:3000/events',
        '/academy': 'http://localhost:3000/main/academy',
        // Discover sub-routes
        '/discover/president-letter': 'http://localhost:3000/discover/president-letter',
        '/discover/board-of-directors': 'http://localhost:3000/discover/board-of-directors',
        '/discover/vision-mission': 'http://localhost:3000/discover/vision-mission',
        '/discover/regulation': 'http://localhost:3000/discover/regulation',
        '/discover/pension-fund': 'http://localhost:3000/discover/pension-fund',
        '/discover/mutual-fund': 'http://localhost:3000/discover/mutual-fund'
    };
    
    // Get the endpoint for current path, default to home
    let endpoint = routeMap[path] || routeMap['/'];
    
    // Append query parameters if they exist (for pagination, filters, etc.)
    if (search) {
        endpoint += search;
    }
    
    console.log('Loading from:', endpoint);
    
    // Use HTMX to load the content
    htmx.ajax('GET', endpoint, {
        target: '#main-div',
        swap: 'innerHTML'
    });
}

// ALWAYS wait for DOMContentLoaded - fixes race condition with cached pages
document.addEventListener('DOMContentLoaded', function() {
    console.log('[LACPA] DOMContentLoaded fired, loading page...');
    loadPageFromURL();
    
    // Keyboard focus styles
    document.querySelectorAll('a').forEach(a => {
        a.addEventListener('focus', () => a.classList.add('ring-2', 'ring-white/25', 'rounded'));
        a.addEventListener('blur', () => a.classList.remove('ring-2', 'ring-white/25', 'rounded'));
    });
});

// Handle browser back/forward buttons
window.addEventListener('popstate', function(event) {
    console.log('[LACPA] Popstate event, loading page...');
    loadPageFromURL();
});

// Anime.js navigation glow animation
let glowAnimation = null;

function initGlowAnimation() {
    const travelingGlow = document.getElementById('traveling-glow');
    const navContainer = travelingGlow?.parentElement;
    const desktopNav = document.querySelector('ul.hidden.md\\:flex');

    if (travelingGlow && navContainer && desktopNav) {
        // Stop existing animation if any
        if (glowAnimation) {
            glowAnimation.pause();
        }

        // Check if we're on desktop (md breakpoint and above)
        const isDesktop = window.innerWidth >= 768;

        if (isDesktop) {
            // Show glow animation only on desktop
            travelingGlow.style.display = 'block';

            // Get responsive glow width based on screen size
            let glowWidth = 60; // mobile default
            if (window.innerWidth >= 1024) {
                glowWidth = 88; // lg screens
            } else if (window.innerWidth >= 768) {
                glowWidth = 70; // md screens
            }

            // Calculate the maximum travel distance
            const navWidth = navContainer.offsetWidth;
            const padding = window.innerWidth >= 768 ? 24 : 12; // responsive padding
            const maxTravel = navWidth - glowWidth - padding;

            // Create the continuous traveling animation
            glowAnimation = anime({
                targets: '#traveling-glow',
                left: [
                    { value: padding / 2, duration: 0 }, // start position
                    { value: maxTravel, duration: 3000 }, // travel to right
                    { value: padding / 2, duration: 3000 } // travel back to left
                ],
                easing: 'easeInOutQuad',
                loop: true,
                direction: 'normal'
            });
        } else {
            // Hide glow animation on mobile
            travelingGlow.style.display = 'none';
        }
    }
}

// Initialize animation on DOM load
document.addEventListener('DOMContentLoaded', initGlowAnimation);

// Reinitialize animation on window resize
let resizeTimeout;
window.addEventListener('resize', function () {
    clearTimeout(resizeTimeout);
    resizeTimeout = setTimeout(initGlowAnimation, 100);
});

// Navigation scroll behavior
function initScrollNavigation() {
    const nav = document.getElementById('main-nav');
    const headerElement = document.querySelector('[hx-get="components/header.html"]'); // The header.html container

    if (!nav) return;

    let headerHeight = 0;
    let isInitialized = false;

    // Function to update header height and position nav
    function updateHeaderHeight() {
        if (headerElement) {
            const newHeight = headerElement.offsetHeight || 0;
            if (newHeight > 0) {
                headerHeight = newHeight;
                isInitialized = true;
                // Update nav position immediately if in static mode
                if (nav.classList.contains('nav-static')) {
                    nav.style.top = `${headerHeight + 20}px`; // 20px offset below header
                }
            }
        }
    }

    // Initial height calculation (likely to be 0 on first load)
    updateHeaderHeight();

    // Update height after HTMX loads the header - THIS IS KEY
    document.addEventListener('htmx:afterSwap', function (event) {
        if (event.detail.target === headerElement) {
            setTimeout(() => {
                updateHeaderHeight();
                // Force initial position update after header loads
                handleScroll();
            }, 100);
        }
    });

    // Also listen for htmx:afterSettle as backup
    document.addEventListener('htmx:afterSettle', function (event) {
        if (event.detail.target === headerElement && !isInitialized) {
            setTimeout(() => {
                updateHeaderHeight();
                handleScroll();
            }, 100);
        }
    });

    // Scroll handler
    function handleScroll() {
        // Don't run scroll logic until header is loaded
        if (!isInitialized || headerHeight === 0) {
            updateHeaderHeight();
            return;
        }

        const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
        const threshold = headerHeight - 20; // Switch slightly before reaching the header bottom

        if (scrollTop > threshold) {
            // Switch to fixed positioning (viewport based)
            if (nav.classList.contains('nav-static')) {
                nav.classList.remove('nav-static');
                nav.classList.add('nav-fixed');
            }
        } else {
            // Switch to absolute positioning (header relative)
            if (nav.classList.contains('nav-fixed')) {
                nav.classList.remove('nav-fixed');
                nav.classList.add('nav-static');
            }
            // Update position based on header height
            if (nav.classList.contains('nav-static')) {
                nav.style.top = `${headerHeight + 20}px`; // 20px offset below header
            }
        }
    }

    // Add scroll listener with throttling
    let scrollTimeout;
    window.addEventListener('scroll', function () {
        if (scrollTimeout) {
            clearTimeout(scrollTimeout);
        }
        scrollTimeout = setTimeout(handleScroll, 10);
    });

    // Initial check - but only after a brief delay to allow header to potentially load
    setTimeout(handleScroll, 200);
}

// Initialize scroll navigation
document.addEventListener('DOMContentLoaded', initScrollNavigation);

// ========================================
// MOBILE MENU FUNCTIONALITY
// ========================================

// Mobile dropdown toggle function (global for onclick)
window.toggleMobileDropdown = function(event, dropdownId) {
    event.preventDefault();
    event.stopPropagation();
    
    const button = event.currentTarget;
    const dropdown = document.getElementById(dropdownId);
    
    if (!dropdown) return;
    
    // Toggle the dropdown
    const isVisible = dropdown.classList.contains('show');
    
    if (isVisible) {
        dropdown.classList.remove('show');
        button.classList.remove('active');
    } else {
        // Close all other dropdowns first
        document.querySelectorAll('.mobile-dropdown').forEach(d => {
            d.classList.remove('show');
        });
        document.querySelectorAll('.mobile-dropdown-toggle').forEach(b => {
            b.classList.remove('active');
        });
        
        // Open this dropdown
        dropdown.classList.add('show');
        button.classList.add('active');
    }
};

function initializeMobileMenu() {
    const menuBtn = document.getElementById('mobile-menu-button');
    const mobileMenu = document.getElementById('mobile-menu');

    if (!menuBtn || !mobileMenu) {
        console.log('Mobile menu elements not found');
        return;
    }

    // Check if already initialized
    if (menuBtn.hasAttribute('data-initialized')) {
        return;
    }

    // Mark as initialized
    menuBtn.setAttribute('data-initialized', 'true');

    const menuIcon = document.getElementById('menu-icon');
    const closeIcon = document.getElementById('close-icon');

    // Simple toggle function
    function toggleMenu() {
        const isExpanded = menuBtn.getAttribute('aria-expanded') === 'true';

        if (isExpanded) {
            // Close menu
            mobileMenu.classList.add('hidden');
            menuBtn.setAttribute('aria-expanded', 'false');
            menuBtn.setAttribute('aria-label', 'Open main menu');
            if (menuIcon) menuIcon.classList.remove('hidden');
            if (closeIcon) closeIcon.classList.add('hidden');
        } else {
            // Open menu
            mobileMenu.classList.remove('hidden');
            menuBtn.setAttribute('aria-expanded', 'true');
            menuBtn.setAttribute('aria-label', 'Close main menu');
            if (menuIcon) menuIcon.classList.add('hidden');
            if (closeIcon) closeIcon.classList.remove('hidden');
        }
    }

    // Add click event listener
    menuBtn.addEventListener('click', function (e) {
        e.preventDefault();
        e.stopPropagation();
        toggleMenu();
    });

    // Close menu when clicking outside
    document.addEventListener('click', function (event) {
        if (!menuBtn.contains(event.target) && !mobileMenu.contains(event.target)) {
            if (!mobileMenu.classList.contains('hidden')) {
                mobileMenu.classList.add('hidden');
                menuBtn.setAttribute('aria-expanded', 'false');
                menuBtn.setAttribute('aria-label', 'Open main menu');
                if (menuIcon) menuIcon.classList.remove('hidden');
                if (closeIcon) closeIcon.classList.add('hidden');
            }
        }
    });
    
    // Close menu when clicking on any navigation link (including dropdowns)
    document.addEventListener('click', function(e) {
        if (e.target.classList.contains('dropdown-link') || 
            e.target.classList.contains('mobile-dropdown-link') ||
            (e.target.closest('a') && e.target.closest('a').hasAttribute('hx-get'))) {
            // Close mobile menu after a short delay to allow HTMX to process
            setTimeout(() => {
                if (!mobileMenu.classList.contains('hidden')) {
                    mobileMenu.classList.add('hidden');
                    menuBtn.setAttribute('aria-expanded', 'false');
                    menuBtn.setAttribute('aria-label', 'Open main menu');
                    if (menuIcon) menuIcon.classList.remove('hidden');
                    if (closeIcon) closeIcon.classList.add('hidden');
                }
            }, 100);
        }
    });
}

// Initialize when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initializeMobileMenu);
} else {
    initializeMobileMenu();
}

// Handle HTMX loading - this is the key part for dynamic content
document.addEventListener('htmx:afterSwap', function (event) {
    if (event.detail.target.querySelector('#mobile-menu-button') ||
        event.detail.target.id === 'mobile-menu-button') {
        setTimeout(initializeMobileMenu, 100);
    }
});

// Additional HTMX event listeners
document.addEventListener('htmx:afterSettle', function (event) {
    setTimeout(initializeMobileMenu, 100);
});

// ========================================
// MAP INITIALIZATION
// ========================================

// Function to initialize Leaflet map
function initializeMap() {
    // Check if map element exists
    const mapElement = document.getElementById('map');
    if (!mapElement) return;

    // Check if map is already initialized
    if (mapElement.hasAttribute('data-map-initialized')) return;

    // Mark as initialized
    mapElement.setAttribute('data-map-initialized', 'true');

    // Custom location coordinates - 33°56'00.2"N 35°37'58.2"E
    const latitude = 33.9334;   // 33°56'00.2"N converted to decimal
    const longitude = 35.6328;  // 35°37'58.2"E converted to decimal
    const zoomLevel = 15;       // Adjust zoom (1-19, higher = closer)

    // Initialize the map
    const map = L.map('map').setView([latitude, longitude], zoomLevel);

    // Add OpenStreetMap tiles
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        maxZoom: 19,
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map);

    // Add marker for your custom location
    L.marker([latitude, longitude]).addTo(map)
        .bindPopup('<b>Dynamic Eye Technology</b><br >Al Ahmarani<br>Block B – 1st Floor')
        .openPopup();

    // Force map to refresh its size
    setTimeout(function () {
        map.invalidateSize();
    }, 100);
}

// Initialize on DOM load
document.addEventListener('DOMContentLoaded', initializeMap);

// Initialize after HTMX swaps content
document.addEventListener('htmx:afterSwap', function(event) {
    // Wait a brief moment for the DOM to settle, then initialize map
    setTimeout(initializeMap, 100);
});

// Initialize after HTMX settles content (backup)
document.addEventListener('htmx:afterSettle', function(event) {
    setTimeout(initializeMap, 100);
});

// ========================================
// SPLIDE CAROUSEL INITIALIZATION
// ========================================

// Function to initialize all Splide carousels
function initializeSplideCarousels() {
    // Hero Carousel - Auto-playing with infinite loop
    const heroCarousel = document.getElementById('hero-carousel');
    if (heroCarousel && !heroCarousel.hasAttribute('data-splide-initialized')) {
        heroCarousel.setAttribute('data-splide-initialized', 'true');
        new Splide('#hero-carousel', {
            type: 'loop',
            autoplay: true,
            interval: 5000,
            pauseOnHover: false,
            pauseOnFocus: false,
            resetProgress: false,
            arrows: true,
            pagination: true,
            height: '62rem',
            width: '100%',
            speed: 1000,
            easing: 'ease-in-out',
            fixedWidth: false,
            fixedHeight: false,
        }).mount();
    }

    // Event Cards Carousel
    const eventCarousel = document.getElementById('splide2');
    if (eventCarousel && !eventCarousel.hasAttribute('data-splide-initialized')) {
        eventCarousel.setAttribute('data-splide-initialized', 'true');
        new Splide('#splide2', {
            type: 'loop',
            perPage: 4,
            perMove: 1,
            gap: '1rem',
            pagination: false,
            arrows: false,
            breakpoints: {
                640: {
                    perPage: 1.5,  // Show 1.5 cards on phone screens
                    gap: '0.75rem', // Slightly smaller gap on mobile
                },
                768: {
                    perPage: 2,    // Show 2 cards on tablet screens
                },
                1024: {
                    perPage: 3,    // Show 3 cards on smaller desktop screens
                },
            },
        }).mount();
    }

    // Publication Cards Carousel
    const publicationCarousel = document.getElementById('splide');
    if (publicationCarousel && !publicationCarousel.hasAttribute('data-splide-initialized')) {
        publicationCarousel.setAttribute('data-splide-initialized', 'true');
        new Splide('#splide', {
            type: 'loop',
            perPage: 5,  // Increased from 3 to 5 to make cards smaller
            perMove: 1,
            gap: '0.75rem',  // Slightly reduced gap
            pagination: false,
            arrows: false,
            breakpoints: {
                640: {
                    perPage: 2,    // Increased from 1.5 to 2 cards on phone screens
                    gap: '0.5rem', // Reduced gap on mobile
                },
                768: {
                    perPage: 3,    // Increased from 2 to 3 cards on tablet screens
                    gap: '0.75rem',
                },
                1024: {
                    perPage: 4,    // Increased from 3 to 4 cards on smaller desktop screens
                    gap: '0.75rem',
                },
                1280: {
                    perPage: 5,    // Increased from 4 to 5 cards on medium desktop screens
                    gap: '0.75rem',
                },
            },
        }).mount();
    }
}

// Initialize on DOM load
document.addEventListener('DOMContentLoaded', initializeSplideCarousels);

// Initialize after HTMX swaps content
document.addEventListener('htmx:afterSwap', function(event) {
    // Wait a brief moment for the DOM to settle, then initialize carousels
    setTimeout(initializeSplideCarousels, 100);
});

// Initialize after HTMX settles content (backup)
document.addEventListener('htmx:afterSettle', function(event) {
    setTimeout(initializeSplideCarousels, 100);
});

// ========================================
// MEMBER CARD NAVIGATION
// ========================================

// Member Card Navigation - Initialize after HTMX swap or DOMContentLoaded
function initializeMemberCards() {
    // Initialize all member cards
    const cards = document.querySelectorAll('.member-card');

    cards.forEach(card => {
        // Skip if already initialized to prevent duplicate event listeners
        if (card.dataset.initialized === 'true') return;
        card.dataset.initialized = 'true';

        let currentState = 0;
        const maxStates = 2; // 0, 1, 2 = 3 states total

        // Get all state elements
        const states = [
            card.querySelector('.card-state-0'),
            card.querySelector('.card-state-1'),
            card.querySelector('.card-state-2')
        ];

        // Function to update card state with animation
        function updateCardState(newState) {
            if (newState < 0 || newState > maxStates) return;

            const currentStateEl = states[currentState];
            const newStateEl = states[newState];

            // Fade out current state
            currentStateEl.classList.add('fade-out');
            currentStateEl.classList.remove('fade-in');

            // Use Anime.js for smooth transition
            anime({
                targets: currentStateEl,
                opacity: 0,
                scale: 0.95,
                duration: 300,
                easing: 'easeInOutQuad',
                complete: function () {
                    currentStateEl.style.pointerEvents = 'none';
                    currentStateEl.style.position = 'absolute';

                    // Fade in new state
                    newStateEl.style.position = 'relative';
                    newStateEl.style.pointerEvents = 'auto';
                    newStateEl.classList.remove('fade-out');
                    newStateEl.classList.add('fade-in');

                    anime({
                        targets: newStateEl,
                        opacity: [0, 1],
                        scale: [0.95, 1],
                        duration: 300,
                        easing: 'easeOutQuad'
                    });
                }
            });

            currentState = newState;
            card.dataset.state = currentState;
        }

        // Event delegation for navigation buttons
        card.addEventListener('click', function (e) {
            const target = e.target.closest('.nav-arrow, .share-btn, .close-share');
            if (!target) return;

            e.preventDefault();
            e.stopPropagation();

            // Add click animation for buttons
            anime({
                targets: target,
                scale: [1, 0.9, 1],
                duration: 200,
                easing: 'easeInOutQuad'
            });

            // Handle share button
            if (target.classList.contains('share-btn')) {
                const shareOverlay = card.querySelector('.share-overlay');
                shareOverlay.style.pointerEvents = 'auto';
                anime({
                    targets: shareOverlay,
                    opacity: [0, 1],
                    scale: [0.95, 1],
                    duration: 400,
                    easing: 'easeOutQuad'
                });
            }
            // Handle close share button
            else if (target.classList.contains('close-share')) {
                const shareOverlay = card.querySelector('.share-overlay');
                anime({
                    targets: shareOverlay,
                    opacity: 0,
                    scale: 0.95,
                    duration: 300,
                    easing: 'easeInQuad',
                    complete: function () {
                        shareOverlay.style.pointerEvents = 'none';
                    }
                });
            }
            // Handle navigation
            else if (target.classList.contains('nav-forward')) {
                // Move forward
                if (currentState < maxStates) {
                    updateCardState(currentState + 1);
                }
            } else if (target.classList.contains('nav-back')) {
                // Move backward
                if (currentState > 0) {
                    updateCardState(currentState - 1);
                }
            }
        });

        // Initialize first state
        states[0].style.position = 'relative';
        states[0].style.pointerEvents = 'auto';
    });
}

// Initialize on DOMContentLoaded (for direct page loads)
document.addEventListener('DOMContentLoaded', function() {
    initializeMemberCards();
});

// Re-initialize after HTMX swaps content
document.body.addEventListener('htmx:afterSwap', function(event) {
    // Check if the swapped content contains member cards
    const target = event.detail.target;
    if (target.querySelector('.member-card') || target.classList.contains('member-card')) {
        initializeMemberCards();
    }
});

// Also handle htmx:load event for additional safety
document.body.addEventListener('htmx:load', function(event) {
    const element = event.detail.elt;
    if (element.querySelector('.member-card') || element.classList.contains('member-card')) {
        initializeMemberCards();
    }
});

// ========================================
// EVENT PAGE TAB FILTERING
// ========================================

// Event category tab filtering
function initializeEventTabs() {
    const tabButtons = document.querySelectorAll('.tab-button');
    const eventCards = document.querySelectorAll('.event-card');

    if (tabButtons.length === 0 || eventCards.length === 0) return;

    tabButtons.forEach(button => {
        button.addEventListener('click', function() {
            const category = this.getAttribute('data-category');

            // Update active tab
            tabButtons.forEach(btn => {
                btn.classList.remove('active', 'text-white');
                btn.classList.add('text-slate-400');
            });
            this.classList.add('active', 'text-white');
            this.classList.remove('text-slate-400');

            // Filter event cards with animation
            eventCards.forEach(card => {
                const cardCategory = card.getAttribute('data-category');
                
                if (category === 'all' || cardCategory === category) {
                    // Show card with fade-in animation
                    card.style.display = 'block';
                    anime({
                        targets: card,
                        opacity: [0, 1],
                        scale: [0.95, 1],
                        duration: 400,
                        easing: 'easeOutQuad'
                    });
                } else {
                    // Hide card with fade-out animation
                    anime({
                        targets: card,
                        opacity: 0,
                        scale: 0.95,
                        duration: 300,
                        easing: 'easeInQuad',
                        complete: function() {
                            card.style.display = 'none';
                        }
                    });
                }
            });
        });
    });
}

// Initialize on DOMContentLoaded
document.addEventListener('DOMContentLoaded', function() {
    initializeEventTabs();
});

// Re-initialize after HTMX swaps content
document.body.addEventListener('htmx:afterSwap', function(event) {
    const target = event.detail.target;
    if (target.querySelector('.tab-button') || target.querySelector('.event-card')) {
        initializeEventTabs();
    }
    
    // Scroll to top of events section after pagination
    if (target.id === 'main-div') {
        setTimeout(() => {
            target.scrollIntoView({ behavior: 'smooth', block: 'start' });
        }, 100);
    }
});

// End of IIFE - execute immediately
})();
