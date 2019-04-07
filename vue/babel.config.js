module.exports = {
    presets: ['@vue/app'],

    plugins: [
        [
            'prismjs',
            {
                languages: ['javascript', 'css', 'markup', 'go'],
                plugins: ['line-numbers'],
                theme: 'twilight',
                css: true,
            },
        ],
    ],
};
