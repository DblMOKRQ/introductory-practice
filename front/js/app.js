// Глобальные переменные
let authHeader = null;
let currentReservationVin = null;

// Авторизация
function handleAuth() {
    if (authHeader) {
        logout();
    } else {
        showAuthModal();
    }
}

// Модальные окна
function showAuthModal() {
    document.getElementById('authModal').style.display = 'flex';
}

function closeAuthModal() {
    document.getElementById('authModal').style.display = 'none';
}

// Авторизация администратора
async function handleLogin(e) {
    e.preventDefault();
    const username = document.getElementById('authLogin').value;
    const password = document.getElementById('authPassword').value;
    authHeader = 'Basic ' + btoa(`${username}:${password}`);
    
    try {
        const response = await fetch('http://localhost:8080/admin/check', {
            headers: { Authorization: authHeader }
        });
        
        if (!response.ok) throw new Error('Неверные данные');
        
        closeAuthModal();
        updateAuthUI();
        showSuccess('Успешный вход');
    } catch (error) {
        showError(error.message);
        authHeader = null;
    }
}

// Выход из системы
function logout() {
    authHeader = null;
    updateAuthUI();
    showSuccess('Сессия завершена');
}

// Обновление UI авторизации
function updateAuthUI() {
    const authButton = document.getElementById('authButton');
    if (authHeader) {
        authButton.textContent = 'Выйти';
        authButton.style.background = 'var(--error)';
    } else {
        authButton.textContent = 'Войти как администратор';
        authButton.style.background = 'var(--secondary)';
    }
}
async function updateStatus(vin, newStatus) {
    if (!authHeader) {
        showError('Требуется аутентификация администратора');
        return;
    }

    showLoading();
    try {
        const response = await fetch(`http://localhost:8080/updateStatus/${vin}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': authHeader
            },
            body: JSON.stringify({ status: newStatus })
        });

        if (response.status === 401) {
            logout();
            throw new Error('Сессия истекла');
        }

        if (!response.ok) throw new Error('Ошибка обновления статуса');
        
        await loadVehicles();
        showSuccess('Статус обновлен');
    } catch (error) {
        showError(error.message);
    } finally {
        hideLoading();
    }
}


// Загрузка автомобилей
document.addEventListener('DOMContentLoaded', loadVehicles);

async function loadVehicles() {
    showLoading();
    try {
        const response = await fetch('http://localhost:8080/getall');
        const vehicles = await response.json();
        renderVehicles(vehicles);
    } catch (error) {
        showError(error.message);
    } finally {
        hideLoading();
    }
}

// Отрисовка автомобилей
function renderVehicles(vehicles) {
    const container = document.getElementById('vehicleList');
    container.innerHTML = vehicles.map(vehicle => `
        <div class="card">
            <h3>${vehicle.brand} ${vehicle.model}</h3>
            <p style="margin: 10px 0; color: #666;">VIN: ${vehicle.vin}</p>
            <p>Год выпуска: ${vehicle.year}</p>
            
            <div class="status-control">
                <select class="status-select" 
                        onchange="updateStatus('${vehicle.vin}', this.value)"
                        value="${vehicle.status}">
                    <option value="available" ${vehicle.status === 'available' ? 'selected' : ''}>Доступен</option>
                    <option value="on_route" ${vehicle.status === 'on_route' ? 'selected' : ''}>В рейсе</option>
                    <option value="under_maintenance" ${vehicle.status === 'under_maintenance' ? 'selected' : ''}>На ТО</option>
                </select>
                <div class="status-indicator ${vehicle.status}"></div>
            </div>

            <div style="margin-top: 15px; display: flex; gap: 10px;">
                <button class="button button-success" 
                    onclick="handleReservationAttempt('${vehicle.vin}', '${vehicle.status}')">
                    <i class="mdi mdi-car-key"></i> Забронировать
                </button>
                <button class="button button-danger" 
                    onclick="checkAuthBeforeDelete('${vehicle.vin}')">
                    <i class="mdi mdi-delete"></i> Удалить
                </button>
            </div>
        </div>
    `).join('');
}

function handleReservationAttempt(vin, currentStatus) {
    if (currentStatus !== 'available') {
        Swal.fire({
            icon: 'error',
            title: 'Невозможно забронировать',
            text: 'Автомобиль недоступен для бронирования. Текущий статус: ' + 
                 getStatusLabel(currentStatus)
        });
        return;
    }
    showReserveModal(vin);
}
function checkAuthBeforeDelete(vin) {
    if (!authHeader) {
        showError('Для удаления требуется войти как администратор');
        return;
    }
    deleteVehicle(vin);
}
function getStatusLabel(status) {
    const labels = {
        available: 'Доступен',
        on_route: 'В рейсе',
        under_maintenance: 'На обслуживании',
    };
    return labels[status] || status;
}

        function showReserveModal(vin) {
            currentReservationVin = vin;
            document.getElementById('reserveModal').style.display = 'flex';
        }

        function closeReserveModal() {
            document.getElementById('reserveModal').style.display = 'none';
            currentReservationVin = null;
        }

        async function handleReservation(e) {
            
            e.preventDefault();
            
            // Получаем текущий статус автомобиля
            const vin = currentReservationVin;
            const vehicle = await getVehicle(vin);
            
            if (vehicle.status !== 'available') {
                showError('Автомобиль больше недоступен для бронирования');
                closeReserveModal();
                await loadVehicles();
                return;
            }
            const userData = {
                name: document.getElementById('userName').value,
                email: document.getElementById('userEmail').value,
                phone: document.getElementById('userPhone').value,
                description: document.getElementById('userDescription').value
            };

            showLoading();
            try {
                const response = await fetch(`http://localhost:8080/rent/${currentReservationVin}`, {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify(userData)
                });

                if (!response.ok) throw new Error('Ошибка бронирования');
                
                closeReserveModal();
                await loadVehicles();
                showSuccess('Автомобиль забронирован');
            } catch (error) {
                showError(error.message);
            } finally {
                hideLoading();
            }
        }

        
        async function deleteVehicle(vin) {
    if (!authHeader) {
        showError('Требуется аутентификация администратора');
        return;
    }

    const { isConfirmed } = await Swal.fire({
        title: 'Удалить ТС?',
        text: 'Это действие нельзя отменить!',
        icon: 'warning',
        showCancelButton: true
    });

    if (isConfirmed) {
        showLoading();
        try {
            const response = await fetch(`http://localhost:8080/delete/${vin}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': authHeader
                }
            });

            if (response.status === 401) {
                logout();
                throw new Error('Сессия истекла');
            }

            if (!response.ok) throw new Error('Ошибка удаления');
            
            await loadVehicles();
            showSuccess('ТС удалено');
        } catch (error) {
            showError(error.message);
        } finally {
            hideLoading();
        }
    }
}

        function showModal(type) {
            document.getElementById('addModal').style.display = 'flex';
        }

        function closeModal() {
            document.getElementById('addModal').style.display = 'none';
            document.getElementById('reserveModal').style.display = 'none';
        }

        async function handleSubmit(e) {
            e.preventDefault();
            const vehicle = {
                vin: document.getElementById('vin').value,
                brand: document.getElementById('brand').value,
                model: document.getElementById('model').value,
                year: parseInt(document.getElementById('year').value),
                status: document.getElementById('status').value
            };

            showLoading();
            try {
                const response = await fetch('http://localhost:8080/add', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify(vehicle)
                });

                if (!response.ok) throw new Error('Ошибка добавления');
                
                closeModal();
                await loadVehicles();
                showSuccess('ТС добавлено');
            } catch (error) {
                showError(error.message);
            } finally {
                hideLoading();
            }
        }

        function showLoading() {
            document.getElementById('loading').style.display = 'flex';
        }

        function hideLoading() {
            document.getElementById('loading').style.display = 'none';
        }

        function showSuccess(message) {
            Swal.fire({ icon: 'success', title: message, timer: 2000 });
        }

        function showError(message) {
            Swal.fire({ icon: 'error', title: 'Ошибка', text: message });

        }
        async function getVehicle(vin) {
    try {
        const response = await fetch(`http://localhost:8080/get/${vin}`);
        return await response.json();
    } catch (error) {
        showError('Ошибка получения данных автомобиля');
        return null;
    }
}