import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const walls = el.querySelector('path[d^="M12.35 21"]');
  animate(walls, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const door = el.querySelector('path[d^="M14.8 12.4"]');
  animate(door, { scaleY: [0, 1.25, 1] }, { duration: τmin, delay: τmid });

  const plus = el.querySelectorAll('path[d="M15 18h6"], path[d="M18 15v6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
