import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const screen = el.querySelector('rect[y="2"]');
  animate(screen, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const light = el.querySelector('path[d="M6 18h.01"]');
  animate(light, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
