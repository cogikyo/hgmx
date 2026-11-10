import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const woofer = el.querySelectorAll('circle, path[d="M12 14h.01"]');
  animate(woofer, { scale: [1, 1.2, 0.9, 1] }, { duration: τmax });

  const tweeter = el.querySelector('path[d="M12 6h.01"]');
  animate(tweeter, { scale: [1, 1.75, 1] }, { duration: τmin, delay: τmid });
}
