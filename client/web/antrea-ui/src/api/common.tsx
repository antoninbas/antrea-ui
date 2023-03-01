export class APIError extends Error {
    code: number;
    status: string;
    date: Date;

    constructor(code: number, status: string, ...params: any[]) {
        super(...params);

        // Maintains proper stack trace for where our error was thrown (only available on V8)
        if (Error.captureStackTrace) {
            Error.captureStackTrace(this, APIError);
        }

        this.name = 'APIError';
        this.code = code;
        this.status = status;
        this.date = new Date();
        this.message = `${this.message} (${this.code}, ${this.status})`
    }
}

export const handleErrorResponse = async (response: Response) => {
    if (!response.ok) {
        console.error("Returned error code: " + response.status);
        console.error("Returned error status:" + response.statusText);
        let errorMessage = "Error processing request.";

        try {
            const errorResponse = await response.json();
            if (errorResponse) {
                errorMessage = errorResponse;
            }
        } catch (e) {
            // do nothing
        } finally {
            throw new APIError(response.status, response.statusText, errorMessage);
        }
    }
    return response.json();
};
