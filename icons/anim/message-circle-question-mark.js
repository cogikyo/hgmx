import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bubble = el.querySelector('path[d^="M2.992"]');
  animate(bubble, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const mark = el.querySelectorAll('path[d^="M9.09"], path[d^="M12 17"]');
  animate(mark, { rotateY: [0, 360], scale: [1, 1.35, 1] }, { duration: τmid, delay: τmin });
}
