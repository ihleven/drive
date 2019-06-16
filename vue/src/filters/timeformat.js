const timeformat = function(ts) {
    if (!ts) return '';
    if (typeof ts == 'string') ts = new Date(ts);
    console.log(typeof ts);
    const today = new Date();

    if (ts.getFullYear() == today.getFullYear()) {
        if (ts.getMonth() == today.getMonth()) {
            if (ts.getDate() == today.getDate()) {
                // 15:30:34
                return ts.getHours() + ':' + ts.getMinutes() + ':' + ts.getSeconds();
            }
            // 12. Jun. 16h
            return ts.getDate() + '. ' + ts.toLocaleString('de', { month: 'short' }) + '. ' + ts.getHours() + 'h';
        }
        // 12. Jun
        return ts.getDate() + '. ' + ts.toLocaleString('de', { month: 'short' }) + '.';
    } else {
        // return ts.toLocaleDateString('de');
        // 2018, Dez.
        return ts.getFullYear() + ', ' + ts.toLocaleString('de', { month: 'short' }) + '.';
    }
};

export default timeformat;
