import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const shield = el.querySelector('path[d^="M20 13"]');
  animate(shield, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const mark = el.querySelectorAll('path[d^="M9.1 9"], path[d="M12 17h.01"]');
  animate(mark, { rotateY: [0, 360], scale: [1, 1.35, 1] }, { duration: τmid, delay: τmin });
}
