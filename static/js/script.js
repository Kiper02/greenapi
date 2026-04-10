const API_BASE_URL = 'http://localhost:8080/api';

function getAuthParams() {
    return {
        idInstance: document.getElementById('idInstance').value,
        apiTokenInstance: document.getElementById('apiTokenInstance').value
    };
}

async function callAPI(method) {
    const { idInstance, apiTokenInstance } = getAuthParams();
    if (!idInstance || !apiTokenInstance) {
        alert('Пожалуйста, заполните idInstance и apiTokenInstance');
        return;
    }

    const url = `${API_BASE_URL}/${method}?idInstance=${idInstance}&apiTokenInstance=${apiTokenInstance}`;
    const resultElem = document.getElementById(`${method}Result`);
    resultElem.textContent = 'Загрузка...';

    try {
        const response = await fetch(url);
        if(!response.ok) {
            resultElem.textContent = await response.text()
            return
        }
        const data = await response.json();
        resultElem.textContent = JSON.stringify(data, null, 2);
    } catch (error) {
        console.log(error)
        resultElem.textContent = `Ошибка: ${error.message}`;
    }
}

async function callSendMessage() {
    const { idInstance, apiTokenInstance } = getAuthParams();
    if (!idInstance || !apiTokenInstance) {
        alert('Пожалуйста, заполните idInstance и apiTokenInstance');
        return;
    }

    const chatId = document.getElementById('chatIdMessage').value;
    const message = document.getElementById('messageText').value;
    if (!chatId || !message) {
        alert('Заполните chatId и текст сообщения');
        return;
    }

    const url = `${API_BASE_URL}/sendMessage?idInstance=${idInstance}&apiTokenInstance=${apiTokenInstance}`;
    const resultElem = document.getElementById('sendMessageResult');
    resultElem.textContent = 'Отправка...';

    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ chatId, message })
        });
        if(!response.ok) {
            resultElem.textContent = await response.text()
            return
        }
        const data = await response.json();
        resultElem.textContent = JSON.stringify(data, null, 2);
    } catch (error) {
        resultElem.textContent = `Ошибка: ${error.message}`;
    }
}

async function callSendFileByUrl() {
    const { idInstance, apiTokenInstance } = getAuthParams();
    if (!idInstance || !apiTokenInstance) {
        alert('Пожалуйста, заполните idInstance и apiTokenInstance');
        return;
    }

    const chatId = document.getElementById('chatIdFile').value;
    const urlFile = document.getElementById('fileUrl').value;
    const fileName = document.getElementById('fileName').value;
    const caption = document.getElementById('caption').value;

    if (!chatId || !urlFile) {
        alert('Заполните chatId и URL файла');
        return;
    }

    const apiUrl = `${API_BASE_URL}/sendFileByUrl?idInstance=${idInstance}&apiTokenInstance=${apiTokenInstance}`;
    const resultElem = document.getElementById('sendFileByUrlResult');
    resultElem.textContent = 'Отправка...';

    try {
        const response = await fetch(apiUrl, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ chatId, urlFile, fileName, caption })
        });
        const data = await response.json();
        resultElem.textContent = JSON.stringify(data, null, 2);
    } catch (error) {
        resultElem.textContent = `Ошибка: ${error.message}`;
    }
}