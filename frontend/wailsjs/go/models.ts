export namespace service {
	
	export class Config {
	    displayName: string;
	    serviceName: string;
	    defaultPort: number;
	    category: string;
	    dependencies: string[];
	    healthCheckType: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.displayName = source["displayName"];
	        this.serviceName = source["serviceName"];
	        this.defaultPort = source["defaultPort"];
	        this.category = source["category"];
	        this.dependencies = source["dependencies"];
	        this.healthCheckType = source["healthCheckType"];
	    }
	}
	export class LogMessage {
	    timestamp: string;
	    level: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new LogMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timestamp = source["timestamp"];
	        this.level = source["level"];
	        this.message = source["message"];
	    }
	}
	export class Status {
	    name: string;
	    isRunning: boolean;
	    message: string;
	    port?: number;
	    healthy?: boolean;
	    category?: string;
	
	    static createFrom(source: any = {}) {
	        return new Status(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.isRunning = source["isRunning"];
	        this.message = source["message"];
	        this.port = source["port"];
	        this.healthy = source["healthy"];
	        this.category = source["category"];
	    }
	}

}

