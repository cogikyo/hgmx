import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const bowl = el.querySelector('path[d^="M12 17"]');
  animate(bowl, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const scoops = el.querySelectorAll('path[d^="M12.14"], path[d^="M15.5"]');
  animate(scoops, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
