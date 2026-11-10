import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const glass = el.querySelector("g");
  animate(glass, { rotate: [0, 195, 180], scale: [1, 0.85, 1] }, { duration: τmax });
  const caps = el.querySelectorAll('path[d$="h14"]');
  animate(caps, { scaleX: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
