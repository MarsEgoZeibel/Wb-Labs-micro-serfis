document.addEventListener('DOMContentLoaded', function() {
    const orderInfo = document.querySelector('.order-info');
    const jsonView = document.getElementById('jsonView');
    const orderID = window.location.pathname.split('/').pop();

    function loadOrderDetails() {
        fetch(`/api/orders/${orderID}`)
            .then(response => response.json())
            .then(order => {
                displayOrderDetails(order);
                displayJSON(order);
            })
            .catch(error => console.error('Error:', error));
    }

    function displayOrderDetails(order) {
        orderInfo.innerHTML = `
            <div class="info-section">
                <h3>Основная информация</h3>
                <div class="info-grid">
                    <div>ID заказа: ${order.order_uid}</div>
                    <div>Трек номер: ${order.track_number}</div>
                    <div>Сервис доставки: ${order.delivery_service}</div>
                    <div>Дата создания: ${new Date(order.date_created).toLocaleString()}</div>
                </div>
            </div>

            <div class="info-section">
                <h3>Информация о доставке</h3>
                <div class="info-grid">
                    <div>Получатель: ${order.delivery.name}</div>
                    <div>Телефон: ${order.delivery.phone}</div>
                    <div>Email: ${order.delivery.email}</div>
                    <div>Город: ${order.delivery.city}</div>
                    <div>Адрес: ${order.delivery.address}</div>
                </div>
            </div>

            <div class="info-section">
                <h3>Оплата</h3>
                <div class="info-grid">
                    <div>Транзакция: ${order.payment.transaction}</div>
                    <div>Сумма: ${order.payment.amount}</div>
                    <div>Валюта: ${order.payment.currency}</div>
                    <div>Банк: ${order.payment.bank}</div>
                </div>
            </div>

            <div class="info-section">
                <h3>Товары</h3>
                ${order.items.map(item => `
                    <div class="info-grid" style="margin-bottom: 10px; padding: 10px; background: #f8f9fa; border-radius: 4px;">
                        <div>Название: ${item.name}</div>
                        <div>Бренд: ${item.brand}</div>
                        <div>Цена: ${item.price}</div>
                        <div>Итого: ${item.total_price}</div>
                        <div>Статус: ${item.status}</div>
                    </div>
                `).join('')}
            </div>
        `;
    }

    function displayJSON(order) {
        jsonView.textContent = JSON.stringify(order, null, 2);
    }
    loadOrderDetails();
});