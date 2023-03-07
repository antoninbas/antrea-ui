const apiServer = process.env.API_SERVER || "";
const apiVersion = "v1";

const config = {
    apiServer: apiServer,
    apiVersion: apiVersion,
    apiUri: `${apiServer}/api/${apiVersion}`,
};

export default config;
