import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const front = el.querySelector('path[d^="M17.5 12"]');
  const back = el.querySelector('path[d^="M21.832"]');
  animate([front, back], { translateX: [-3, 1, 0] }, { duration: τmax });
  animate(front, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });
  animate(back, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
