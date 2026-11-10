import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const apple = el.querySelector('path[d^="M18.237"]');
  animate(apple, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const stem = el.querySelector('path[d^="M12 6.528"]');
  animate(stem, { rotate: [0, -25, 10, 0] }, { duration: τmid, delay: τmin });
}
