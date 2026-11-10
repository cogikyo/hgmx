import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const arc = el.querySelector('path[d*="10 10 0"]');
  animate(arc, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const arrow = el.querySelectorAll('path:not([d*="10 10 0"])');
  animate(arrow, { translateX: [-4, 1, 0], translateY: [4, -1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
