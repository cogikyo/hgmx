import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const note = el.querySelector('path[d^="M12 18"]');
  animate(note, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const arrow = el.querySelectorAll('path[d^="m16 19"], path[d="M19 16v6"]');
  animate(arrow, { translateY: [-3, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
