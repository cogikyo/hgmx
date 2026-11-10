import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const post = el.querySelectorAll('path[d^="M12"]');
  animate(post, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const sign = el.querySelector('path[d^="M2.354"]');
  animate(sign, { translateX: [0, -3, 3, -1, 0] }, { duration: τmax });
}
