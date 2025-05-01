export const handler = async(event) => {
    const response = {
        statusCode: 200,
        headers: {
            "X-Some-Header": "Wow"
        },
        body: JSON.stringify('Hello from Lambda!'),
    };
    return response;
};
