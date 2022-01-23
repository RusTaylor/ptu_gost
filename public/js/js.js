$('#sendForm').on('click', (el) => {
    $.ajax({
        type: "POST",
        url: "/saveperson",
        data: {
            fio: $('#fio').val(),
            birthday: $('#birthday').val(),
            code: $('#code').val(),
        },

        success: function () {
            $('#fio').val(null)
            $('#birthday').val(null)
            $('#code').val(null)
            $('#sendResult').removeClass('none')
            setTimeout(() => {
                $('#sendResult').addClass('none')
            }, 5000)
        }
    });
})