import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const dot = el.querySelector('path[d="M12 12h.01"]');
  animate([el.querySelector("g"), dot], { translateX: [0, -4, 2, 0], translateY: [0, -3, -2, 0] }, { duration: τmax });
  animate(dot, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
