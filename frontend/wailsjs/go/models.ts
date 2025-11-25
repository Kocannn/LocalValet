export namespace domain {
	
	export class ServiceStatus {
	    name: string;
	    isRunning: boolean;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new ServiceStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.isRunning = source["isRunning"];
	        this.message = source["message"];
	    }
	}

}

export namespace main {
	
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
	export class ServiceConfig {
	    DisplayName: string;
	    Linux: string;
	    Darwin: string;
	    Windows: string;
	
	    static createFrom(source: any = {}) {
	        return new ServiceConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.DisplayName = source["DisplayName"];
	        this.Linux = source["Linux"];
	        this.Darwin = source["Darwin"];
	        this.Windows = source["Windows"];
	    }
	}

}

