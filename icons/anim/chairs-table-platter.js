import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const chairs = el.querySelectorAll('path[d^="M8 19"], path[d^="M2 8"], path[d^="M16 19"], path[d^="M22 8"]');
  animate(chairs, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const cloche = el.querySelectorAll('path[d^="M8 10"], path[d^="M12 6"]');
  animate(cloche, { translateY: [0, -3, 0] }, { duration: τmid, delay: τmin });
}
