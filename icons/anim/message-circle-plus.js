import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bubble = el.querySelector('path[d^="M2.992"]');
  animate(bubble, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const plus = el.querySelectorAll('path[d="M8 12h8"], path[d="M12 8v8"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
