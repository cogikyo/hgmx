import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelectorAll('path[d^="M12.543"], path[d^="M7.606"]');
  animate(body, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const plus = el.querySelectorAll('path[d="M10 9v6"], path[d="M7 12h6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
