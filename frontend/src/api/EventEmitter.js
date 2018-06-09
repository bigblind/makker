export default class EventEmitter {
    constructor() {
        this.listeners = new Map();
    }

    on(event, handler) {
        this.listeners.has(event) || this.listeners.set(event, []);
        this.listeners.get(event).push(handler);
    }

    off(event, handler){
        let listeners = this.listeners.get(event);

        if (listeners && listeners.length) {
            for (let i=0; i<listeners.length; i++) {
                if (listeners[i] === handler) {
                    listeners.splice(i, 1);
                    this.listeners.set(event, listeners);
                    return true
                }
            }
        }

        return false;
    }

    emit(event, ...args) {
        let listeners = this.listeners.get(event)
        if (listeners) {
            listeners.forEach((l) => {
                console.log(l)
                l(...args);
            })
        }
    }
}