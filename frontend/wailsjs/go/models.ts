export namespace domain {
	
	export class BinarySourceInfo {
	    os: string;
	    using_system_binaries: boolean;
	    binary_location: string;
	    binary_validation: any;
	
	    static createFrom(source: any = {}) {
	        return new BinarySourceInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.os = source["os"];
	        this.using_system_binaries = source["using_system_binaries"];
	        this.binary_location = source["binary_location"];
	        this.binary_validation = source["binary_validation"];
	    }
	}
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

}

