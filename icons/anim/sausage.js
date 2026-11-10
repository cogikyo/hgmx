import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const sausage = el.querySelector('path[d^="M22 19"]');
  animate(sausage, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const fork = el.querySelector("g");
  animate(fork, { translateX: [0, 1.5, 0], translateY: [0, -1.5, 0] }, { duration: τmid, delay: τmin });
}
