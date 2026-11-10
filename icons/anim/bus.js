import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const bus = el.querySelector("g");
  const wheels = el.querySelectorAll("circle");
  animate([bus, ...wheels], { translateX: [-6, 1, 0], opacity: [0, 1] }, { duration: τmax });
  animate(bus, { rotate: [0, 3, -1, 0] }, { duration: τmid, delay: τmin });
}
