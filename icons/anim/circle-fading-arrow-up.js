import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const arc = el.querySelector('path[d^="M12 2a"]');
  animate(arc, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const arrow = el.querySelectorAll('path[d^="m16 12"], path[d^="M12 16"]');
  animate(arrow, { translateY: [4, -1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
