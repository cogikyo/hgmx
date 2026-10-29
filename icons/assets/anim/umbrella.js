import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const canopy = el.querySelector('path[d^="M20.992"]');
  animate(canopy, { scaleX: [0.25, 1.15, 1] }, { duration: τmid });

  const handle = el.querySelector('path[d^="M12 13"]');
  animate(handle, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
