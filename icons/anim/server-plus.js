import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const racks = el.querySelectorAll('path[d^="M12.5"], path[d^="M22 18"]');
  animate(racks, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: 0.3, delay: stagger(0.09) });

  const plus = el.querySelectorAll('path[d="M16 12h6"], path[d="M19 9v6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
