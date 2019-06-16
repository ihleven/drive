const ext = ['B', 'kB', 'MB', 'GB', 'TB'];

const bytes = function(size) {
    if (!size) return '';

    let j = 0;

    while (size > 1000) {
        size = size / 1000;
        j++;
    }

    return +size.toFixed(1) + ext[j];
};

export default bytes;

