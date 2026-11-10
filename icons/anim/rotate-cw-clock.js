import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const hands = el.querySelector('path[d^="M12 7"]');
  animate(hands, { rotate: [0, 360] }, { duration: τmax });

  const head = el.querySelector('path[d^="M16 8"]');
  animate(head, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
