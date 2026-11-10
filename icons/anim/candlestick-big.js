import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const candle = el.querySelector('path[d^="M16 22"]');
  animate(candle, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const wick = el.querySelector('path[d="M12 11v3"]');
  animate(wick, { scaleY: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
