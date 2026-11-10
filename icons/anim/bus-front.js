import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const bus = el.querySelector("g");
  animate(bus, { scale: [0.5, 1.1, 1], opacity: [0, 1] }, { duration: τmax });
  animate(bus, { rotate: [0, 3, -1, 0] }, { duration: τmid, delay: τmin });
}
