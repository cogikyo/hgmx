import { animate, τmax } from "../motion.js";

export default function (el) {
  const long = el.querySelector('rect[height="16"]');
  animate(long, { translateY: [-1.5, 0.5, 0] }, { duration: τmax });

  const short = el.querySelector('rect[height="9"]');
  animate(short, { translateY: [-3.5, 0.5, 0] }, { duration: τmax });
}
