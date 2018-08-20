export function buildQueryString(parameters) {
    const keys = Object.keys(parameters);
    if (keys.length === 0) {
        return '';
    }

    let query = '?';
    for (let i = 0; i < keys.length; i++) {
        const key = keys[i];
        query += key + '=' + encodeURIComponent(parameters[key]);

        if (i < keys.length - 1) {
            query += '&';
        }
    }

    return query;
}