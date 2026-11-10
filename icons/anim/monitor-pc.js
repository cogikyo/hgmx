import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const screen = el.querySelector('path[d^="M10 15"]');
  animate(screen, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const light = el.querySelector('path[d="M18 17h.01"]');
  animate(light, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: τmin });
}
