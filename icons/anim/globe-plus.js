import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const globe = el.querySelectorAll('path[d^="M15.929"], path[d^="M21.92"], path[d="M22 12H2"]');
  animate(globe, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const plus = el.querySelectorAll('path[d="M19 2v6"], path[d="M22 5h-6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
