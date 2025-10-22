/**
 * CMS Admin Dashboard JavaScript
 * Uses SweetAlert2 for all notifications, confirmations, and alerts
 * Theme: Dark mode (#1f1f1f background, #ffffff text)
 */

// Initialize mobile sidebar functionality
function initializeMobileSidebar() {
    const mobileMenuBtn = document.getElementById('mobile-menu-btn');
    const mobileSidebar = document.getElementById('mobile-sidebar');
    const mobileOverlay = document.getElementById('mobile-sidebar-overlay');
    
    if (!mobileMenuBtn || !mobileSidebar || !mobileOverlay) {
        console.log('Mobile sidebar elements not found');
        return;
    }
    
    console.log('Initializing mobile sidebar...');
    
    function toggleMobileSidebar() {
        mobileSidebar.classList.toggle('-translate-x-full');
        mobileOverlay.classList.toggle('hidden');
        document.body.classList.toggle('overflow-hidden'); // Prevent background scroll
    }
    
    // Mobile menu button click
    mobileMenuBtn.addEventListener('click', function(e) {
        e.preventDefault();
        e.stopPropagation();
        console.log('Mobile menu button clicked');
        toggleMobileSidebar();
    });
    
    // Overlay click to close
    mobileOverlay.addEventListener('click', function(e) {
        e.preventDefault();
        console.log('Mobile overlay clicked');
        toggleMobileSidebar();
    });
    
    // Close mobile sidebar when clicking a link
    const mobileNavLinks = document.querySelectorAll('.mobile-nav-link');
    mobileNavLinks.forEach(link => {
        link.addEventListener('click', function() {
            console.log('Mobile nav link clicked');
            toggleMobileSidebar();
        });
    });
}

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
        
        // Adjust main content margin (only on large screens)
        if (mainContent) {
            mainContent.classList.remove('lg:ml-64');
            mainContent.classList.add('lg:ml-20');
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
        
        // Adjust main content margin (only on large screens)
        if (mainContent) {
            mainContent.classList.remove('lg:ml-20');
            mainContent.classList.add('lg:ml-64');
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
    
    // Check if the swapped element is the sidebar (desktop or mobile)
    if (event.detail.target.querySelector('#sidebar') || event.detail.target.id === 'sidebar' ||
        event.detail.target.querySelector('#mobile-sidebar') || event.detail.target.id === 'mobile-sidebar') {
        console.log('Sidebar detected in swap');
        setTimeout(initializeSidebar, 50);
        setTimeout(initializeMobileSidebar, 50);
    }
    
    // Check if header was swapped (contains mobile menu button)
    if (event.detail.target.querySelector('#mobile-menu-btn') || event.detail.target.id === 'mobile-menu-btn') {
        console.log('Header detected in swap');
        setTimeout(initializeMobileSidebar, 50);
    }
    
    // Check if slide content was swapped
    if (event.detail.target.id === 'section-content') {
        console.log('Slide content loaded');
        // Update active tab styling based on which button triggered the request
        const triggerElement = event.detail.elt;
        if (triggerElement && triggerElement.classList.contains('slide-tab')) {
            setActiveSlideTab(triggerElement);
        }

        // Initialize drag-and-drop for upload areas
        initializeUploadArea();
    }
});

// Also initialize on page load in case sidebar is already loaded
document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM loaded, waiting for sidebar...');
    // Wait a bit for HTMX to load components
    setTimeout(function() {
        initializeSidebar();
        initializeMobileSidebar();
    }, 200);
    
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
        const response = await fetch('/api/admin/slides', {
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
    // Remove active state from all tabs
    document.querySelectorAll('.slide-tab').forEach(tab => {
        tab.classList.remove('text-white', 'border-b-2', 'border-blue-500');
        tab.classList.add('text-gray-500');
    });
    
    // Add active state to clicked tab
    activeTab.classList.remove('text-gray-500');
    activeTab.classList.add('text-white', 'border-b-2', 'border-blue-500');
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
        
        // Reload tabs from server via HTMX instead of manually updating
        const slideTabsWrapper = document.getElementById('slide-tabs-wrapper');
        if (slideTabsWrapper && typeof htmx !== 'undefined') {
            // Trigger HTMX to reload the tabs
            htmx.ajax('GET', 'http://localhost:3000/api/admin/slides/tabs', {
                target: '#slide-tabs-wrapper',
                swap: 'innerHTML'
            }).then(() => {
                // After tabs are reloaded, load the new slide content
                setTimeout(() => {
                    const newSlideTab = slideTabsWrapper.querySelector(`.slide-tab[hx-get*="${newSlide.id}"]`);
                    if (newSlideTab) {
                        htmx.trigger(newSlideTab, 'click');
                    }
                }, 100);
            });
        }
        
        showNotification(`New slide added successfully!`, 'success');
    } catch (error) {
        console.error('Error creating slide:', error);
        showNotification('Error creating slide', 'error');
    }
}

async function deleteSlide(slideId) {
    // Confirm deletion with SweetAlert2
    const result = await Swal.fire({
        title: 'Delete Slide?',
        text: 'This action cannot be undone!',
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#dc2626',
        cancelButtonColor: '#6b7280',
        confirmButtonText: 'Yes, delete it!',
        cancelButtonText: 'Cancel',
        background: '#1f1f1f',
        color: '#ffffff',
        customClass: {
            popup: 'border border-gray-700'
        }
    });

    if (!result.isConfirmed) {
        return;
    }

    // Show loading state
    Swal.fire({
        title: 'Deleting...',
        text: 'Please wait',
        icon: 'info',
        allowOutsideClick: false,
        showConfirmButton: false,
        background: '#1f1f1f',
        color: '#ffffff',
        didOpen: () => {
            Swal.showLoading();
        }
    });

    try {
        const response = await fetch(`http://localhost:3000/api/admin/slides/${slideId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`
            }
        });

        if (!response.ok) {
            console.error('Failed to delete slide');
            Swal.fire({
                title: 'Error!',
                text: 'Failed to delete slide',
                icon: 'error',
                confirmButtonColor: '#3b82f6',
                background: '#1f1f1f',
                color: '#ffffff'
            });
            return;
        }

        // Reload tabs from server via HTMX
        const slideTabsWrapper = document.getElementById('slide-tabs-wrapper');
        if (slideTabsWrapper && typeof htmx !== 'undefined') {
            htmx.ajax('GET', 'http://localhost:3000/api/admin/slides/tabs', {
                target: '#slide-tabs-wrapper',
                swap: 'innerHTML'
            }).then(() => {
                // After tabs are reloaded, load the first slide if available
                setTimeout(() => {
                    const firstTab = slideTabsWrapper.querySelector('.slide-tab');
                    if (firstTab) {
                        htmx.trigger(firstTab, 'click');
                    } else {
                        // No slides left, clear the content area
                        document.getElementById('section-content').innerHTML = '<div class="text-center text-gray-400 py-12">No slides available. Click "Add Slide" to create one.</div>';
                    }
                }, 100);
            });
        }

        Swal.fire({
            title: 'Deleted!',
            text: 'Slide deleted successfully',
            icon: 'success',
            timer: 2000,
            showConfirmButton: false,
            background: '#1f1f1f',
            color: '#ffffff'
        });
    } catch (error) {
        console.error('Error deleting slide:', error);
        Swal.fire({
            title: 'Error!',
            text: 'Error deleting slide',
            icon: 'error',
            confirmButtonColor: '#3b82f6',
            background: '#1f1f1f',
            color: '#ffffff'
        });
    }
}

async function toggleSlideField(checkbox) {
    const slideId = checkbox.dataset.slideId;
    const field = checkbox.dataset.field;
    const isChecked = checkbox.checked;

    // Show loading toast
    const Toast = Swal.mixin({
        toast: true,
        position: 'top-end',
        showConfirmButton: false,
        timer: 1500,
        background: '#1f1f1f',
        color: '#ffffff'
    });

    Toast.fire({
        icon: 'info',
        title: 'Updating...'
    });

    try {
        const response = await fetch(`http://localhost:3000/api/admin/slides/${slideId}`, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`
            },
            body: JSON.stringify({
                [field]: isChecked
            })
        });

        if (!response.ok) {
            console.error('Failed to update toggle');
            // Revert checkbox state
            checkbox.checked = !isChecked;
            
            Swal.fire({
                title: 'Error!',
                text: `Failed to update ${field.replace('Active', '')}`,
                icon: 'error',
                confirmButtonColor: '#3b82f6',
                background: '#1f1f1f',
                color: '#ffffff',
                timer: 2000
            });
            return;
        }

        // Success feedback
        Toast.fire({
            icon: 'success',
            title: `${field.replace('Active', '')} ${isChecked ? 'enabled' : 'disabled'}`
        });

    } catch (error) {
        console.error('Error updating toggle:', error);
        // Revert checkbox state
        checkbox.checked = !isChecked;
        
        Swal.fire({
            title: 'Error!',
            text: 'Network error occurred',
            icon: 'error',
            confirmButtonColor: '#3b82f6',
            background: '#1f1f1f',
            color: '#ffffff',
            timer: 2000
        });
    }
}

async function saveSlideChanges(slideId) {
    // Get form values
    const buttonTitle = document.getElementById(`buttonTitle-${slideId}`)?.value || '';
    const buttonLink = document.getElementById(`buttonLink-${slideId}`)?.value || '';
    const title = document.getElementById(`title-${slideId}`)?.value || '';
    const description = document.getElementById(`description-${slideId}`)?.value || '';

    // Validate required fields
    if (!title.trim()) {
        Swal.fire({
            title: 'Validation Error',
            text: 'Title is required',
            icon: 'warning',
            confirmButtonColor: '#3b82f6',
            background: '#1f1f1f',
            color: '#ffffff'
        });
        return;
    }

    // Show loading state
    Swal.fire({
        title: 'Saving Changes...',
        text: 'Please wait',
        icon: 'info',
        allowOutsideClick: false,
        showConfirmButton: false,
        background: '#1f1f1f',
        color: '#ffffff',
        didOpen: () => {
            Swal.showLoading();
        }
    });

    try {
        const response = await fetch(`http://localhost:3000/api/admin/slides/${slideId}`, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`
            },
            body: JSON.stringify({
                title: title,
                description: description,
                buttonTitle: buttonTitle,
                buttonLink: buttonLink
            })
        });

        if (!response.ok) {
            throw new Error('Failed to save changes');
        }

        const updatedSlide = await response.json();

        // Success feedback
        Swal.fire({
            title: 'Saved!',
            text: 'Changes saved successfully',
            icon: 'success',
            timer: 2000,
            showConfirmButton: false,
            background: '#1f1f1f',
            color: '#ffffff'
        });

    } catch (error) {
        console.error('Error saving changes:', error);
        Swal.fire({
            title: 'Error!',
            text: 'Failed to save changes. Please try again.',
            icon: 'error',
            confirmButtonColor: '#3b82f6',
            background: '#1f1f1f',
            color: '#ffffff'
        });
    }
}

async function cancelSlideChanges(slideId) {
    // Confirm cancel action
    const result = await Swal.fire({
        title: 'Discard Changes?',
        text: 'Any unsaved changes will be lost',
        icon: 'question',
        showCancelButton: true,
        confirmButtonColor: '#6b7280',
        cancelButtonColor: '#3b82f6',
        confirmButtonText: 'Yes, discard',
        cancelButtonText: 'Keep editing',
        background: '#1f1f1f',
        color: '#ffffff'
    });

    if (!result.isConfirmed) {
        return;
    }

    // Reload the slide from server to reset form
    try {
        const slideContent = document.getElementById('section-content');
        if (slideContent && typeof htmx !== 'undefined') {
            htmx.ajax('GET', `http://localhost:3000/api/admin/slides/${slideId}/render`, {
                target: '#section-content',
                swap: 'innerHTML'
            });

            const Toast = Swal.mixin({
                toast: true,
                position: 'top-end',
                showConfirmButton: false,
                timer: 2000,
                background: '#6b7280',
                color: '#ffffff'
            });

            Toast.fire({
                icon: 'info',
                title: 'Changes discarded'
            });
        }
    } catch (error) {
        console.error('Error reloading slide:', error);
    }
}

async function handleSlideImageUpload(input) {
    const slideId = input.dataset.slideId;
    const file = input.files[0];

    if (!file) {
        return;
    }

    // Validate file type
    if (!file.type.startsWith('image/')) {
        Swal.fire({
            title: 'Invalid File',
            text: 'Please upload an image file',
            icon: 'error',
            confirmButtonColor: '#3b82f6',
            background: '#1f1f1f',
            color: '#ffffff'
        });
        input.value = ''; // Clear input
        return;
    }

    // Validate file size (10MB)
    if (file.size > 10 * 1024 * 1024) {
        Swal.fire({
            title: 'File Too Large',
            text: 'Image size must be less than 10MB',
            icon: 'error',
            confirmButtonColor: '#3b82f6',
            background: '#1f1f1f',
            color: '#ffffff'
        });
        input.value = ''; // Clear input
        return;
    }

    // Show loading state
    Swal.fire({
        title: 'Uploading Image...',
        text: 'Please wait',
        icon: 'info',
        allowOutsideClick: false,
        showConfirmButton: false,
        background: '#1f1f1f',
        color: '#ffffff',
        didOpen: () => {
            Swal.showLoading();
        }
    });

    try {
        // Create FormData
        const formData = new FormData();
        formData.append('image', file);

        // Upload image
        const response = await fetch(`http://localhost:3000/api/admin/slides/${slideId}/upload-image`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`
            },
            body: formData
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Failed to upload image');
        }

        const result = await response.json();

        // Update the current image preview
        const currentImageContainer = document.getElementById(`current-image-${slideId}`);
        if (currentImageContainer) {
            currentImageContainer.innerHTML = `
                <img 
                    src="${result.url}?t=${Date.now()}" 
                    alt="Current section image" 
                    class="w-full h-full object-cover">
                <div class="absolute inset-0 bg-gradient-to-t from-black/50 to-transparent"></div>
            `;
        }

        // Success feedback
        Swal.fire({
            title: 'Success!',
            text: 'Image uploaded successfully',
            icon: 'success',
            timer: 2000,
            showConfirmButton: false,
            background: '#1f1f1f',
            color: '#ffffff'
        });

        // Clear input
        input.value = '';

    } catch (error) {
        console.error('Error uploading image:', error);
        Swal.fire({
            title: 'Upload Failed',
            text: error.message || 'Failed to upload image. Please try again.',
            icon: 'error',
            confirmButtonColor: '#3b82f6',
            background: '#1f1f1f',
            color: '#ffffff'
        });
        // Clear input
        input.value = '';
    }
}

function initializeUploadArea() {
    // Find all upload areas in the current slide
    const uploadAreas = document.querySelectorAll('[id^="upload-area-"]');
    
    uploadAreas.forEach(uploadArea => {
        const slideId = uploadArea.id.replace('upload-area-', '');
        const fileInput = document.getElementById(`file-input-${slideId}`);
        
        if (!fileInput) return;

        // Click to upload
        uploadArea.addEventListener('click', (e) => {
            // Prevent triggering when clicking the button
            if (!e.target.closest('button')) {
                fileInput.click();
            }
        });

        // Drag over
        uploadArea.addEventListener('dragover', (e) => {
            e.preventDefault();
            uploadArea.classList.add('drag-over');
        });

        // Drag leave
        uploadArea.addEventListener('dragleave', (e) => {
            e.preventDefault();
            uploadArea.classList.remove('drag-over');
        });

        // Drop
        uploadArea.addEventListener('drop', (e) => {
            e.preventDefault();
            uploadArea.classList.remove('drag-over');
            
            const files = e.dataTransfer.files;
            if (files.length > 0) {
                // Set the file to the input and trigger the change event
                const dataTransfer = new DataTransfer();
                dataTransfer.items.add(files[0]);
                fileInput.files = dataTransfer.files;
                
                // Trigger the upload
                handleSlideImageUpload(fileInput);
            }
        });
    });
}

function initializeTabSwitching() {
    const slideTabsContainer = document.getElementById('slide-tabs');
    
    // Handle section tab clicks
    document.addEventListener('click', function(e) {
        const sectionTab = e.target.closest('.section-tab');
        if (sectionTab) {
            // Remove active state from all section tabs
            document.querySelectorAll('.section-tab').forEach(tab => {
                tab.classList.remove('text-white', 'border-b-2', 'border-blue-500', '-mb-px');
                tab.classList.add('text-gray-400');
            });
            
            // Add active state to clicked tab
            sectionTab.classList.remove('text-gray-400');
            sectionTab.classList.add('text-white', 'border-b-2', 'border-blue-500', '-mb-px');
            
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

        // Handle slide tab clicks - update active state immediately
        const slideTab = e.target.closest('.slide-tab');
        if (slideTab) {
            console.log('Slide tab clicked:', slideTab);
            setActiveSlideTab(slideTab);
        }
    });
}

function showNotification(message, type = 'success') {
    const Toast = Swal.mixin({
        toast: true,
        position: 'top-end',
        showConfirmButton: false,
        timer: 3000,
        timerProgressBar: true,
        background: type === 'success' ? '#059669' : '#dc2626',
        color: '#ffffff',
        didOpen: (toast) => {
            toast.addEventListener('mouseenter', Swal.stopTimer);
            toast.addEventListener('mouseleave', Swal.resumeTimer);
        }
    });

    Toast.fire({
        icon: type,
        title: message
    });
}


