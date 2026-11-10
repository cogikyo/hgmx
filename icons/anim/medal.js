import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const medal = el.querySelector("g");
  animate(medal, { rotate: [0, -10, 6, -3, 0] }, { duration: τmax });
  const one = el.querySelector('path[d^="M12 18"]');
  animate(one, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
