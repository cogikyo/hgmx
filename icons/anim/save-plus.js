import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const disk = el.querySelectorAll('path[d^="M12.5"], path[d^="M16 13"]');
  animate(disk, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const plus = el.querySelectorAll('path[d="M19 22v-6"], path[d="M22 19h-6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
