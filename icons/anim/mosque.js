import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const walls = el.querySelectorAll('path[d^="M7 22V6"], path[d="M3 9h4"]');
  animate(walls, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const crescent = el.querySelector('path[d^="M12.268"]');
  animate(crescent, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
