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
