import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const walls = el.querySelectorAll('path[d^="M18 11"], path[d^="M18 21"]');
  animate(walls, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const cross = el.querySelectorAll('path[d^="M12 7"], path[d^="M14 9"]');
  animate(cross, { scale: [1, 0.85, 1.15, 1] }, { duration: τmid, delay: τmin });
}
