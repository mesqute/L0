// ждем пока загрузится страница
$(document).ready(function () {
    $('#id-btn').click(function () {
        $.ajax({
            type: "GET",
            url: "data",
            data: {
                'id': $('#id-input').val()
            },
            dataType: 'json',
            error : function (_, _, error) {
                alert(error, "danger")
            },
            success: function (response) {
                response.ModelState
                // очищаем таблицу order
                $('#order-row').empty();
                // обходим все JSON поля ответа
                $.each(response, function (_, field) {
                    // если поле delivery
                    if (field === response["delivery"]) {
                        // очищаем таблицу delivery
                        $('#delivery-row').empty();
                        $.each(field, function (idx, cell) {
                            addCell('#delivery-row', cell)
                        });
                        // переходим на следующее поле (continue)
                        return;
                    }
                    // если поле payment
                    if (field === response["payment"]) {
                        // очищаем таблицу payment
                        $('#payment-row').empty();
                        $.each(field, function (idx, cell) {
                            addCell('#payment-row', cell)
                        });
                        // переходим на следующее поле (continue)
                        return;
                    }
                    // если поле items
                    if (field === response["items"]) {
                        let row = '#items-row'
                        // очищаем таблицу items
                        $('#items-body').empty();
                        $.each(response["items"], function (index, elem) {
                            // генерируем номерной id для новой строки таблицы
                            let rowId = row + index
                            // вставляем новую строку
                            $('#items-body').append($("<tr id='items-row" + index + "'></tr>"));
                            // вставляем значния в строку
                            $.each(elem, function (idx, cell) {
                                addCell(rowId, cell)
                            });
                        });
                        // переходим на следующее поле (continue)
                        return;
                    }
                    // если поле не одна из подструктур, то значит это элементы таблицы order
                    addCell('#order-row', field);
                })

            }
        })
    })

});
// генерирует ячейки таблицы в определенную строку
function addCell(row, elem) {
    $(row).append($('<td>', {
        text: elem
    }));

}

// выводит сообщение об ошибке
function alert(message, type) {
    let wrapper = document.createElement('div')
    wrapper.innerHTML = '<div style="z-index: 11" class="alert alert-' + type + ' alert-dismissible" role="alert">' + message + '<button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button></div>'

    $('#AlertPlaceholder').append(wrapper)
}
