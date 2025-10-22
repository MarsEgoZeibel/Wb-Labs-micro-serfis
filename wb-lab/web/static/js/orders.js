let allOrders = [];
const ordersList = document.getElementById('ordersList');
const searchInput = document.getElementById('searchInput');
const searchBtn = document.getElementById('searchBtn');
const resetBtn = document.getElementById('resetBtn');

function formatDate(dateStr) {
    const date = new Date(dateStr);
    return date.toLocaleString('ru-RU', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    });
}

function displayOrders(orders) {
    if (!ordersList) return;
    
    ordersList.innerHTML = '';
    
    if (!orders || orders.length === 0) {
        ordersList.innerHTML = `
            <div style="text-align: center; padding: 40px; color: var(--primary);">
                <h3>🔍 Заказы не найдены</h3>
                <p>Попробуйте изменить параметры поиска</p>
            </div>
        `;
        return;
    }

    orders.forEach(function(order) {
        const card = document.createElement('div');
        card.className = 'order-card';
        
        card.innerHTML = `
            <h3>📦 Заказ #${order.order_uid}</h3>
            <div class="info-row">
                <span class="emoji">🚚</span>
                <span class="label">Трек-номер:</span>
                ${order.track_number}
            </div>
            <div class="info-row">
                <span class="emoji">👤</span>
                <span class="label">Получатель:</span>
                ${order.delivery.name}
            </div>
            <div class="info-row">
                <span class="emoji">📍</span>
                <span class="label">Адрес:</span>
                ${order.delivery.city}, ${order.delivery.address}
            </div>
            <div class="info-row">
                <span class="emoji">💰</span>
                <span class="label">Стоимость:</span>
                ${order.payment.amount} ${order.payment.currency}
            </div>
            <div class="info-row">
                <span class="emoji">🏦</span>
                <span class="label">Банк:</span>
                ${order.payment.bank}
            </div>
            <div class="info-row">
                <span class="emoji">📅</span>
                <span class="label">Дата:</span>
                ${formatDate(order.date_created)}
            </div>
            <div class="status-badge">
                ${order.delivery_service} ✨
            </div>
        `;
        
        card.addEventListener('click', function() {
            window.location.href = '/order/' + order.order_uid;
        });
        
        ordersList.appendChild(card);
    });
}

function loadOrders() {
    fetch('/api/orders')
        .then(function(response) {
            return response.json();
        })
        .then(function(orders) {
            allOrders = orders;
            displayOrders(orders);
        })
        .catch(function(error) {
            console.error('Error:', error);
            if (ordersList) {
                ordersList.innerHTML = `
                    <div style="text-align: center; padding: 40px; color: #ff4444;">
                        <h3>❌ Ошибка загрузки заказов</h3>
                        <p>Пожалуйста, попробуйте обновить страницу</p>
                    </div>
                `;
            }
        });
}

function searchOrders() {
    const searchTerm = searchInput.value.trim().toLowerCase();
    if (searchTerm === '') {
        displayOrders(allOrders);
        return;
    }
    
    const filteredOrders = allOrders.filter(function(order) {
        return order.order_uid.toLowerCase().includes(searchTerm);
    });
    displayOrders(filteredOrders);
}

window.addEventListener('load', function() {
    loadOrders();
    if (searchBtn) {
        searchBtn.addEventListener('click', searchOrders);
    }

    if (resetBtn) {
        resetBtn.addEventListener('click', function() {
            if (searchInput) {
                searchInput.value = '';
            }
            displayOrders(allOrders);
        });
    }

    if (searchInput) {
        searchInput.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                searchOrders();
            }
        });
    }
});