import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const button = el.querySelector('path[d="M10 2h4"]');
  animate(button, { translateY: [0, 2, 0] }, { duration: τmin });
  const hand = el.querySelector('path[d^="M12 14"]');
  animate(hand, { rotate: [0, 360] }, { duration: τmid, delay: τmin });
  const head = el.querySelector('path[d^="M9 17"]');
  animate(head, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
