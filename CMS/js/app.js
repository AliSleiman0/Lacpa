// Initialize sidebar functionality after HTMX swap
function initializeSidebar() {
    const sidebar = document.getElementById('sidebar');
    const toggleBtn = document.getElementById('sidebar-toggle');
    const mainContent = document.querySelector('main');
    
    if (!sidebar || !toggleBtn) {
        console.log('Sidebar or toggle button not found');
        return;
    }
    
    console.log('Initializing sidebar...');
    
    // Load saved state from localStorage
    const isCollapsed = localStorage.getItem('sidebarCollapsed') === 'true';
    if (isCollapsed) {
        collapseSidebar();
    }
    
    // Add toggle event listener
    toggleBtn.addEventListener('click', function(e) {
        e.preventDefault();
        e.stopPropagation();
        console.log('Toggle clicked');
        
        if (sidebar.classList.contains('w-64')) {
            collapseSidebar();
            localStorage.setItem('sidebarCollapsed', 'true');
        } else {
            expandSidebar();
            localStorage.setItem('sidebarCollapsed', 'false');
        }
    });
    
    function collapseSidebar() {
        console.log('Collapsing sidebar');
        sidebar.classList.remove('w-64');
        sidebar.classList.add('w-20');
        const icon = toggleBtn.querySelector('svg');
        if (icon) icon.style.transform = 'rotate(180deg)';
        
        // Hide all text elements
        const textElements = sidebar.querySelectorAll('.sidebar-text');
        textElements.forEach(el => {
            el.style.opacity = '0';
            el.style.width = '0';
            el.style.overflow = 'hidden';
        });
        
        // Adjust main content margin
        if (mainContent) {
            mainContent.classList.remove('ml-64');
            mainContent.classList.add('ml-20');
        }
    }
    
    function expandSidebar() {
        console.log('Expanding sidebar');
        sidebar.classList.remove('w-20');
        sidebar.classList.add('w-64');
        const icon = toggleBtn.querySelector('svg');
        if (icon) icon.style.transform = 'rotate(0deg)';
        
        // Show all text elements
        const textElements = sidebar.querySelectorAll('.sidebar-text');
        textElements.forEach(el => {
            el.style.opacity = '1';
            el.style.width = 'auto';
            el.style.overflow = 'visible';
        });
        
        // Adjust main content margin
        if (mainContent) {
            mainContent.classList.remove('ml-20');
            mainContent.classList.add('ml-64');
        }
    }
    
    // Active state management for sidebar
    const currentPath = window.location.pathname;
    const navLinks = sidebar.querySelectorAll('nav a');
    
    navLinks.forEach(link => {
        const linkPath = new URL(link.href).pathname;
        
        if (currentPath === linkPath || currentPath.includes(link.dataset.page)) {
            // Remove active state from all links
            navLinks.forEach(l => {
                l.classList.remove('bg-blue-600', 'text-white', 'shadow-lg');
                l.classList.add('text-gray-400', 'hover:bg-[#2a2a2a]');
            });
            
            // Add active state to current link
            link.classList.remove('text-gray-400', 'hover:bg-[#2a2a2a]');
            link.classList.add('bg-blue-600', 'text-white', 'shadow-lg');
        }
    });
}

// Listen for HTMX after swap event
document.addEventListener('htmx:afterSwap', function(event) {
    console.log('HTMX afterSwap event triggered');
    // Check if the swapped element is the sidebar
    if (event.detail.target.querySelector('#sidebar') || event.detail.target.id === 'sidebar') {
        console.log('Sidebar detected in swap');
        setTimeout(initializeSidebar, 50);
    }
    
    // Check if slide content was swapped
    if (event.detail.target.id === 'section-content') {
        console.log('Slide content loaded');
        // Update active tab styling based on which button triggered the request
        const triggerElement = event.detail.elt;
        if (triggerElement && triggerElement.classList.contains('slide-tab')) {
            setActiveSlideTab(triggerElement);
        }
    }
});

// Also initialize on page load in case sidebar is already loaded
document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM loaded, waiting for sidebar...');
    // Wait a bit for HTMX to load components
    setTimeout(initializeSidebar, 200);
    
    // Initialize tab switching
    initializeTabSwitching();
    
    // Show slide tabs on initial load since Hero Section is active by default
    const slideTabsContainer = document.getElementById('slide-tabs');
    if (slideTabsContainer) {
        slideTabsContainer.style.display = 'block';
    }
});

// Tab switching functionality
let allSlides = [];

async function loadSlides() {
    try {
        const response = await fetch('http://localhost:3000/api/admin/slides', {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`
            }
        });
        
        if (!response.ok) {
            console.error('Failed to load slides');
            return;
        }
        
        allSlides = await response.json();
        renderSlideTabs();
        
        // Load first slide if available
        if (allSlides.length > 0) {
            loadSlideContent(allSlides[0].id);
        }
    } catch (error) {
        console.error('Error loading slides:', error);
    }
}

function renderSlideTabs() {
    const slideTabsWrapper = document.getElementById('slide-tabs-wrapper');
    const addSlideBtn = document.getElementById('add-slide-btn');
    
    if (!slideTabsWrapper) {
        console.error('Slide tabs wrapper not found');
        return;
    }
    
    // Remove all existing slide tabs (keep only add button)
    const existingTabs = slideTabsWrapper.querySelectorAll('.slide-tab');
    existingTabs.forEach(tab => tab.remove());
    
    // Create tabs for each slide
    allSlides.forEach((slide, index) => {
        const slideTab = document.createElement('button');
        slideTab.className = 'slide-tab px-4 py-2 text-sm font-medium text-gray-500 hover:text-white transition-colors';
        slideTab.dataset.slideId = slide.id;
        slideTab.textContent = `Slide ${index + 1}`;
        
        // Add HTMX attributes for dynamic loading
        slideTab.setAttribute('hx-get', `/api/admin/slides/${slide.id}/render`);
        slideTab.setAttribute('hx-target', '#section-content');
        slideTab.setAttribute('hx-swap', 'innerHTML');
        
        slideTab.addEventListener('click', () => {
            setActiveSlideTab(slideTab);
        });
        
        slideTabsWrapper.insertBefore(slideTab, addSlideBtn);
    });
    
    // Process HTMX for new elements
    if (typeof htmx !== 'undefined') {
        htmx.process(slideTabsWrapper);
    }
    
    // Set first slide as active and load it
    if (allSlides.length > 0) {
        const firstTab = slideTabsWrapper.querySelector('.slide-tab');
        if (firstTab) {
            setActiveSlideTab(firstTab);
            // Trigger HTMX click to load content
            htmx.trigger(firstTab, 'click');
        }
    }
}

function setActiveSlideTab(activeTab) {
    document.querySelectorAll('.slide-tab').forEach(tab => {
        tab.classList.remove('text-white', 'border-blue-500');
        tab.classList.add('text-gray-500');
    });
    
    activeTab.classList.remove('text-gray-500');
    activeTab.classList.add('text-white', 'border-blue-500');
}

async function loadSlideContent(slideId) {
    try {
        const response = await fetch(`http://localhost:3000/api/admin/slides/${slideId}/render`, {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`
            }
        });
        
        if (!response.ok) {
            console.error('Failed to load slide content');
            return;
        }
        
        const html = await response.text();
        const sectionContent = document.getElementById('section-content');
        if (sectionContent) {
            sectionContent.innerHTML = html;
        }
    } catch (error) {
        console.error('Error loading slide content:', error);
    }
}

async function createNewSlide() {
    try {
        const response = await fetch('http://localhost:3000/api/admin/slides', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`
            },
            body: JSON.stringify({
                title: 'New Slide',
                description: '',
                imgSrc: '',
                buttonTitle: '',
                buttonLink: '',
                isActive: true,
                imageActive: true,
                buttonActive: true,
                titleActive: true,
                descriptionActive: true
            })
        });
        
        if (!response.ok) {
            console.error('Failed to create slide');
            showNotification('Failed to create slide', 'error');
            return;
        }
        
        const newSlide = await response.json();
        allSlides.push(newSlide);
        renderSlideTabs();
        loadSlideContent(newSlide.id);
        showNotification(`Slide ${allSlides.length} added successfully!`, 'success');
    } catch (error) {
        console.error('Error creating slide:', error);
        showNotification('Error creating slide', 'error');
    }
}

function initializeTabSwitching() {
    const slideTabsContainer = document.getElementById('slide-tabs');
    
    // Handle section tab clicks
    document.addEventListener('click', function(e) {
        const sectionTab = e.target.closest('.section-tab');
        if (sectionTab) {
            // Remove active state from all section tabs
            document.querySelectorAll('.section-tab').forEach(tab => {
                tab.classList.remove('text-white', 'border-blue-500');
                tab.classList.add('text-gray-400');
            });
            
            // Add active state to clicked tab
            sectionTab.classList.remove('text-gray-400');
            sectionTab.classList.add('text-white', 'border-blue-500');
            
            // Show/hide slide tabs based on section
            if (sectionTab.dataset.section === 'hero') {
                slideTabsContainer.style.display = 'block';
                // Load slides when hero section is selected
                loadSlides();
            } else {
                slideTabsContainer.style.display = 'none';
            }
        }
        
        // Handle add slide button click
        const addSlideBtn = e.target.closest('#add-slide-btn');
        if (addSlideBtn) {
            createNewSlide();
        }
    });
}

function showNotification(message, type = 'success') {
    const bgColor = type === 'success' ? 'bg-green-600' : 'bg-red-600';
    const notification = document.createElement('div');
    notification.className = `fixed top-20 right-4 ${bgColor} text-white px-6 py-3 rounded-lg shadow-lg z-50 transition-all duration-300`;
    notification.innerHTML = `
        <div class="flex items-center gap-3">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="${type === 'success' ? 'M5 13l4 4L19 7' : 'M6 18L18 6M6 6l12 12'}" />
            </svg>
            <span>${message}</span>
        </div>
    `;
    
    document.body.appendChild(notification);
    
    setTimeout(() => {
        notification.style.opacity = '0';
        setTimeout(() => notification.remove(), 300);
    }, 3000);
}


