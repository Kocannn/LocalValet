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
	export class ModuleConfig {
	    displayName: string;
	    serviceName: string;
	    installed: boolean;
	    enabled: boolean;
	    installCommand: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new ModuleConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.displayName = source["displayName"];
	        this.serviceName = source["serviceName"];
	        this.installed = source["installed"];
	        this.enabled = source["enabled"];
	        this.installCommand = source["installCommand"];
	        this.description = source["description"];
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
	export class UserPreferences {
	    enabledModules: string[];
	
	    static createFrom(source: any = {}) {
	        return new UserPreferences(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabledModules = source["enabledModules"];
	    }
	}

}

